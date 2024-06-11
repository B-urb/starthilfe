
# Starthilfe

![Starthilfe Logo](assets/starthilfe.webp)

`Starthilfe` is a command-line tool designed to streamline the process of setting up new projects or integrating with existing ones by using a template repository for standard files and a repository for reusable GitHub Actions. The tool uses a YAML configuration file to determine actions and files to be integrated into the project.

## Features

- **Initialize Projects**: Sets up a new project or integrates with an existing one using a predefined template repository.
- **Manage GitHub Actions**: Adds and updates GitHub Actions as subtrees, based on the programming language specified in the configuration.
- **Configurable**: Uses a YAML file to manage configurations for easy adaptation and reuse.
- **Logging**: Utilizes structured logging via Go's `log/slog` for clear and actionable logs.

## Requirements

- Go 1.21 or later (due to the use of `log/slog` for structured logging).
- Git installed on your system.

## Installation

Clone the repository and build the tool:

```bash
git clone https://github.com/yourusername/starthilfe.git
cd starthilfe
go build -o starthilfe
```

## Configuration

`Starthilfe` operates based on a YAML configuration file named `starthilfe.yml`. Here’s a sample configuration:

```yaml
language: 'python'
template_repo:
  url: https://github.com/yourusername/template-repo.git
  files: ['releaserc', 'codeowners', 'renovate.json']
  branch: 'main'
  force: true
actions_repo:
  url: https://github.com/yourusername/aaand-action.git
  subtree_path: '.github/workflows/actions/python'
  branch: 'main'
```

- **language**: Programming language of the project.
- **template_repo**: Repository containing template files like `releaserc`, `codeowners`.
- **actions_repo**: Repository containing GitHub Actions.

## Usage

`Starthilfe` supports several commands:

### Initialize Configuration

Generate a default `starthilfe.yml` configuration file.

```bash
./starthilfe init
```

### Add Subtrees

Add GitHub Actions as subtrees to your project as specified in the configuration file.

```bash
./starthilfe add
```

### Update Subtrees

Update the GitHub Actions subtrees to the latest version from the repository.

```bash
./starthilfe update
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you have feedback or proposals for improvements.

## License

Specify the license under which the tool is released, such as MIT, GPL, etc.

## Support

For support, please open an issue in the GitHub repository or contact [your email].

---

### Additional Notes

The README provides a comprehensive overview of what the tool does, how to set it up, and how to use it. It includes basic sections that should be present in most open-source projects, such as installation instructions, usage examples, configuration details, contributing guidelines, and license information. Adjust the repository URLs and contact details as per your actual project details.