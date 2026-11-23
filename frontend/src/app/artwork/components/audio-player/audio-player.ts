import {
  ChangeDetectionStrategy,
  Component,
  computed,
  DestroyRef,
  effect,
  ElementRef,
  inject,
  input,
  signal,
  viewChild,
} from '@angular/core';
import { AudioPlaybackService } from '../../services/audio-playback.service';

@Component({
  selector: 'app-audio-player',
  templateUrl: './audio-player.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AudioPlayer {
  private readonly playback = inject(AudioPlaybackService);
  private readonly id = Symbol('audio-player');

  readonly src = input.required<string>();

  constructor() {
    // Pause when another player becomes the active one.
    effect(() => {
      const active = this.playback.activePlayerId();
      const audio = this.audioRef().nativeElement;
      if (active !== this.id && !audio.paused) {
        audio.pause();
      }
    });

    inject(DestroyRef).onDestroy(() => this.playback.release(this.id));
  }

  protected onPlay() {
    this.playing.set(true);
    this.playback.activate(this.id);
  }

  protected onPause() {
    this.playing.set(false);
    this.playback.release(this.id);
  }

  protected readonly audioRef = viewChild.required<ElementRef<HTMLAudioElement>>('audioEl');

  protected readonly playing = signal(false);
  protected readonly currentTime = signal(0);
  protected readonly duration = signal(0);

  protected readonly progress = computed(() =>
    this.duration() ? (this.currentTime() / this.duration()) * 100 : 0,
  );

  protected togglePlay(event: Event) {
    event.stopPropagation();
    const audio = this.audioRef().nativeElement;
    if (audio.paused) {
      audio.play();
    } else {
      audio.pause();
    }
  }

  protected seek(event: Event) {
    event.stopPropagation();
    if (!this.duration()) return;
    const value = +(event.target as HTMLInputElement).value;
    this.audioRef().nativeElement.currentTime = (value / 100) * this.duration();
  }

  protected onEnded() {
    this.onPause();
    this.currentTime.set(0);
  }

  protected formatTime(seconds: number): string {
    if (!isFinite(seconds)) return '0:00';
    const m = Math.floor(seconds / 60);
    const s = Math.floor(seconds % 60)
      .toString()
      .padStart(2, '0');
    return `${m}:${s}`;
  }
}
