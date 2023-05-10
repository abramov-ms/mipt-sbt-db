import ncache


def main() -> None:
    cache = ncache.CacheManager.get_cache("demoCache")


if __name__ == "__main__":
    main()
