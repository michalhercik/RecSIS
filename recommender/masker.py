class Masker:
    def mask(self, reference: list[str], values: list[str]) -> list[bool]:
        """
        keep order of values
        """
        raise NotImplementedError()

class InMasker(Masker):
    def mask(self, reference: list[str], values: list[str]) -> list[bool]:
        mask = [course in reference for course in values]
        return mask


class NotInMasker(InMasker):
    def mask(self, reference: list[str], values: list[str]) -> list[bool]:
        mask = [not m for m in super().mask(reference, values)]
        return mask
