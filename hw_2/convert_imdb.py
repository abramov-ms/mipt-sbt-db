#!/usr/bin/python

import argparse
import sys
import json

from pathlib import Path


def parse_args():
    parser = argparse.ArgumentParser(
        prog="convert_imdb.py",
        description="Convert IMDB reviews dataset into JSON format",
    )
    parser.add_argument(
        "--in",
        metavar="INPUT_DATASET_DIR",
        help="Directory with IMDB dataset, i.e. train/ or test/",
        required=True,
    )
    parser.add_argument(
        "--out",
        metavar="OUTPUT_DATASET_DIR",
        help="Directory with bunch of converted JSON files",
        required=True,
    )

    args_namespace = parser.parse_args(sys.argv[1:])
    return vars(args_namespace)


def parse_path(path: Path) -> tuple[int, int, str]:
    name = path.name
    idx = int(name[: name.find("_")])
    score = int(name[name.find("_") + 1 : name.find(".")])
    review = path.read_text()
    return idx, score, review


def title_id_from_url(url: str) -> int:
    return int(url[url.rfind("tt") + 2 : url.rfind("/")])


def main():
    args = parse_args()

    in_dir = Path(args["in"])
    out_file = Path(args["out"])
    urls_pos = (in_dir / "urls_pos.txt").read_text().split()
    urls_neg = (in_dir / "urls_neg.txt").read_text().split()

    # title_id => [(review, score, type)]
    reviews_by_title: dict[int, list[dict[str, any]]] = dict()

    for path in (in_dir / "pos").glob("*"):
        idx, score, review = parse_path(path)
        title_id = title_id_from_url(urls_pos[idx])
        reviews_by_title.setdefault(title_id, []).append(
            {"review": review, "score": score, "type": "pos"}
        )

    for path in (in_dir / "neg").glob("*"):
        idx, score, review = parse_path(path)
        title_id = title_id_from_url(urls_pos[idx])
        reviews_by_title.setdefault(title_id, []).append(
            {"review": review, "score": score, "type": "neg"}
        )

    flat_reviews = []
    for title_id, reviews in reviews_by_title.items():
        flat_reviews.append({"id": title_id, "reviews": reviews})

    with out_file.open("w+") as f:
        json.dump(flat_reviews, f, indent=2)


if __name__ == "__main__":
    try:
        main()
    except IOError as e:
        print("IO error happened, check the file paths:", e)
        exit(1)
    except Exception as e:
        print("Unknown error:", e)
        exit(1)
