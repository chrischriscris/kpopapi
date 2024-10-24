#!/usr/bin/env python3

import os
import shutil
import sys
from enum import Enum
from typing import Generator

from PIL import Image


class ImageKind(Enum):
    PORTRAIT = 1
    LANDSCAPE = 2
    NOT_ENOUGH_RES = 3


IMAGE_EXTENSIONS = (".jpg", ".jpeg", ".png")
LANDSCAPE_FOLDER = "landscape"
PORTRAIT_FOLDER = "portrait"

min_width_landscape = 1600
min_height_landscape = 900

min_width_portrait = 900
min_height_portrait = 1600


def get_leaf_nodes(directory) -> Generator[str, None, None]:
    for root, _, files in os.walk(directory):
        for file in files:
            yield os.path.join(root, file)


def classify_image(image: Image) -> ImageKind:
    width, height = image.size

    if width / height >= 1:
        return (
            ImageKind.LANDSCAPE
            if width >= min_width_landscape and height >= min_height_landscape
            else ImageKind.NOT_ENOUGH_RES
        )

    return (
        ImageKind.PORTRAIT
        if width >= min_width_portrait and height >= min_height_portrait
        else ImageKind.NOT_ENOUGH_RES
    )


def is_portrait_image(image: Image) -> bool:
    width, height = image.size
    return (
        height / width >= 1
        and width >= min_width_portrait
        and height >= min_height_portrait
    )


def move_to_folder(file: str, folder: str):
    if not os.path.exists(folder):
        os.mkdir(folder)

    shutil.copy(file, folder)


def landscape_callback(file: str):
    move_to_folder(file, LANDSCAPE_FOLDER)


def portrait_callback(file: str):
    move_to_folder(file, PORTRAIT_FOLDER)


def main():
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <directory>")
        sys.exit(1)

    directory = sys.argv[1]
    for file in get_leaf_nodes(directory):
        # If the file is not an image, skip it
        if not file.endswith(IMAGE_EXTENSIONS):
            continue

        try:
            image = Image.open(file)
            kind = classify_image(image)
            match kind:
                case ImageKind.LANDSCAPE:
                    landscape_callback(file)
                case ImageKind.PORTRAIT:
                    portrait_callback(file)
                case _:
                    pass

        except Exception:
            print(f"Could not open image {file}")
            continue


if __name__ == "__main__":
    main()
