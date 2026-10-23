import { AssetResponse } from '../interfaces/asset-response.interface';
import { Asset } from '../interfaces/asset.interface';

export class AssetMapper {
  static fromResponse(resp: AssetResponse): Asset {
    return { id: resp.id, path: resp.path };
  }
}
