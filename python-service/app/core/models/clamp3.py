"""CLaMP 3 (https://github.com/sanderwood/clamp3), audio <-> multilingual text.

Slim re-implementation of the inference path only (text + audio towers). Audio is
first turned into MERT-v1-95M features (5 s segments, mean over time and layers),
which CLaMP 3's audio encoder then pools into a single 768-d vector.
"""

import numpy as np
import torch
from huggingface_hub import hf_hub_download
from torch import Tensor, nn
from transformers import (
    AutoConfig,
    AutoModel,
    AutoTokenizer,
    BertConfig,
    BertModel,
    Wav2Vec2FeatureExtractor,
)

from app.core.models.base import AudioTextModel

CLAMP3_REPO = "sander-wood/clamp3"
CLAMP3_SHA = "355625cc1c6f73726bbcd0eb9276ac7152d56426"
CLAMP3_WEIGHTS = (
    "weights_clamp3_saas_h_size_768_t_model_FacebookAI_xlm-roberta-base_t_length_128"
    "_a_size_768_a_layers_12_a_length_128_s_size_768_s_layers_12_p_size_64_p_length_512.pth"
)

MERT_REPO = "m-a-p/MERT-v1-95M"
MERT_SHA = "12af15fef9d0ac838c3f475bfbbf26d2060dd4f5"
MERT_SR = 24000
MERT_SEGMENT_SEC = 5

TEXT_REPO = "FacebookAI/xlm-roberta-base"
TEXT_SHA = "e73636d4f797dec63c3081bb6ed5c7b0bb3f2089"

HIDDEN = 768
MAX_TEXT_LENGTH = 128
MAX_AUDIO_LENGTH = 128  # MERT segments (+2 boundary vectors) per encoder pass
MERT_BATCH = 16


class _Clamp3Towers(nn.Module):
    def __init__(self):
        super().__init__()
        self.text_model = AutoModel.from_config(
            AutoConfig.from_pretrained(TEXT_REPO, revision=TEXT_SHA)
        )
        self.text_proj = nn.Linear(self.text_model.config.hidden_size, HIDDEN)
        self.audio_model = BertModel(
            BertConfig(
                vocab_size=1,
                hidden_size=HIDDEN,
                num_hidden_layers=12,
                num_attention_heads=HIDDEN // 64,
                intermediate_size=HIDDEN * 4,
                max_position_embeddings=MAX_AUDIO_LENGTH,
            )
        )
        self.audio_proj = nn.Linear(HIDDEN, HIDDEN)

    @staticmethod
    def _avg_pool(features: Tensor, masks: Tensor) -> Tensor:
        masks = masks.unsqueeze(-1).to(features.dtype)
        return (features * masks).sum(dim=1) / masks.sum(dim=1)

    def text_global(self, ids: Tensor, masks: Tensor) -> Tensor:
        h = self.text_model(ids, attention_mask=masks)["last_hidden_state"]
        return self.text_proj(self._avg_pool(h, masks))

    def audio_global(self, feats: Tensor, masks: Tensor) -> Tensor:
        h = self.audio_model(inputs_embeds=feats, attention_mask=masks)[
            "last_hidden_state"
        ]
        return self.audio_proj(self._avg_pool(h, masks))


class Clamp3(AudioTextModel):
    sample_rate = MERT_SR

    def load(self):
        self._device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

        self._tokenizer = AutoTokenizer.from_pretrained(TEXT_REPO, revision=TEXT_SHA)

        self._mert = (
            AutoModel.from_pretrained(
                MERT_REPO, trust_remote_code=True, revision=MERT_SHA
            )
            .to(self._device)
            .eval()
        )
        self._mert_processor = Wav2Vec2FeatureExtractor.from_pretrained(
            MERT_REPO, revision=MERT_SHA
        )

        path = hf_hub_download(CLAMP3_REPO, CLAMP3_WEIGHTS, revision=CLAMP3_SHA)
        checkpoint = torch.load(path, map_location="cpu", weights_only=True)
        # The checkpoint also holds the symbolic (sheet music / MIDI) tower; unused here.
        state = {
            k: v
            for k, v in checkpoint["model"].items()
            if not k.startswith("symbolic_")
        }
        self._towers = _Clamp3Towers()
        self._towers.load_state_dict(state)
        self._towers.to(self._device).eval()

    # ---- text -------------------------------------------------------------

    def encode_text(self, texts: list[str]) -> Tensor:
        out = []
        for text in texts:
            enc = self._tokenizer(
                text, return_tensors="pt", truncation=True, max_length=MAX_TEXT_LENGTH
            )
            ids = enc["input_ids"].to(self._device)
            masks = enc["attention_mask"].to(self._device)
            with torch.no_grad():
                out.append(self._towers.text_global(ids, masks)[0])
        return self._normalize(torch.stack(out))

    # ---- audio ------------------------------------------------------------

    def encode_audio(self, waveforms: list[np.ndarray]) -> Tensor:
        out = [self._encode_one_audio(w) for w in waveforms]
        return self._normalize(torch.stack(out))

    @torch.no_grad()
    def _mert_segment_features(self, wav: np.ndarray) -> Tensor:
        """[n_segments, 768]: per 5 s segment, mean over time and all 13 layers."""
        values = self._mert_processor(
            wav, sampling_rate=MERT_SR, return_tensors="pt"
        ).input_values[0]
        size = MERT_SR * MERT_SEGMENT_SEC
        segments = [values[i : i + size] for i in range(0, len(values), size)]
        # Reference pipeline drops a trailing segment shorter than 1 s.
        if len(segments) > 1 and len(segments[-1]) < MERT_SR:
            segments.pop()

        feats = []
        full = [s for s in segments if len(s) == size]
        partial = [s for s in segments if len(s) != size]
        for i in range(0, len(full), MERT_BATCH):
            feats.append(self._mert_pass(torch.stack(full[i : i + MERT_BATCH])))
        for s in partial:
            feats.append(self._mert_pass(s.unsqueeze(0)))
        return torch.cat(feats, dim=0)

    def _mert_pass(self, batch: Tensor) -> Tensor:
        hidden = self._mert(
            batch.to(self._device), output_hidden_states=True
        ).hidden_states
        return torch.stack(hidden).mean(dim=(0, 2))  # layers and time -> [B, 768]

    @torch.no_grad()
    def _encode_one_audio(self, wav: np.ndarray) -> Tensor:
        feats = self._mert_segment_features(wav).cpu()
        zero = torch.zeros(1, HIDDEN)
        data = torch.cat((zero, feats, zero), 0)

        # Same chunking as the reference extractor: fixed windows, last one aligned
        # to the end, chunk vectors averaged weighted by their number of new frames.
        n, win = len(data), MAX_AUDIO_LENGTH
        chunks = [data[i : i + win] for i in range(0, n, win)]
        chunks[-1] = data[-win:]
        weights = [win] * (n // win) + ([n % win] if n % win else [])

        vecs = []
        for chunk in chunks:
            length = chunk.size(0)
            padded = torch.cat((chunk, torch.zeros(win - length, HIDDEN)), 0)
            mask = torch.cat((torch.ones(length), torch.zeros(win - length)), 0)
            vecs.append(
                self._towers.audio_global(
                    padded.unsqueeze(0).to(self._device),
                    mask.unsqueeze(0).to(self._device),
                )[0]
            )
        w = torch.tensor(weights, dtype=torch.float32, device=self._device)[:, None]
        return (torch.stack(vecs) * w).sum(0) / w.sum()

    @staticmethod
    def _normalize(x: Tensor) -> Tensor:
        return x / x.norm(p=2, dim=-1, keepdim=True)
