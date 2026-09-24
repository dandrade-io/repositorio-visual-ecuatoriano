class ModelRegistry[T]:
    def __init__(self, model_classes: dict[str, type[T]]) -> None:
        self._model_classes: dict[str, type[T]] = model_classes
        self._loaded: dict[str, T] = {}

    def get(self, name: str):
        if name not in self._model_classes:
            raise ValueError(f"Unknown model: {name}")
        if name not in self._loaded:
            self._loaded[name] = self._instantiate(name)

        return self._loaded[name]

    def _instantiate(self, name: str) -> T:
        model_cls = self._model_classes[name]
        instance = model_cls()
        load = getattr(instance, "load", None)

        if callable(load):
            load()

        return instance
