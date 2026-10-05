export enum Predominance {
  Tipography = 'typography',
  Image = 'image',
  Color = 'color',
  ImageAndColor = 'image_and_color',
}

export const PREDOMINANCE_LABELS: Record<Predominance, string> = {
  [Predominance.Image]: 'Imagen',
  [Predominance.Tipography]: 'Tipografía',
  [Predominance.Color]: 'Color',
  [Predominance.ImageAndColor]: 'Imagen y color',
};

