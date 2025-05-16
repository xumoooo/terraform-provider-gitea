# Examples

The examples in this folder serve multiple purposes. Most importantly:
- To provide examples used by the tfplugindocs cli tool to generate documentation (in the `docs/` folder). [These filepaths](https://github.com/hashicorp/terraform-plugin-docs/blob/main/README.md#conventional-paths) are used in this way.
- To provide a central example module (in the root of `examples/`) that in turn uses the valid example modules in `resources/`. This allows the use of the example module as a smoketest for this provider.
