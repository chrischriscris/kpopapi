#!/usr/bin/env python3

import shutil
import sys
import os
from typing import Generator
from PIL import Image

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


def is_landscape_image(image: Image) -> bool:
    width, height = image.size
    return width / height >= 1 and width >= min_width_landscape and height >= min_height_landscape


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
            if is_landscape_image(image):
                landscape_callback(file)
            else:
                pass
                # portrait_callback(file)
        except Exception:
            print(f"Could not open image {file}")
            continue


if __name__ == "__main__":
    main()
