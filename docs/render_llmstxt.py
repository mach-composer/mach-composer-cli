def extract_files_from_nav(nav_item):
    """Recursively extract all file paths from the navigation structure"""
    files = []

    if isinstance(nav_item, str):
        # Direct file path
        files.append(nav_item)
    elif isinstance(nav_item, list):
        # List of nav items
        for item in nav_item:
            files.extend(extract_files_from_nav(item))
    elif isinstance(nav_item, dict):
        # Dictionary with keys and values
        for key, value in nav_item.items():
            files.extend(extract_files_from_nav(value))

    return files


def generate_llm_txt(nav, custom_text=None):
    """Generate llm.txt file from the navigation structure"""
    import os

    # Extract all files from nav
    all_files = extract_files_from_nav(nav)

    # Separate docs from changelog
    docs_files = []
    optional_files = []

    for file_path in all_files:
        if file_path == 'changelog.md':
            optional_files.append(file_path)
        else:
            docs_files.append(file_path)

    # Generate markdown content
    content = "# Mach Composer\n\n"

    # Add custom text if provided
    if custom_text:
        content += f"{custom_text}\n\n"

    # Add docs section
    content += "## Docs\n\n"
    for file_path in docs_files:
        content += f"- {file_path}\n"

    # Add optional section if there are optional files
    if optional_files:
        content += "\n## Optional\n\n"
        for file_path in optional_files:
            content += f"- {file_path}\n"

    # Write to build folder
    build_dir = os.path.join(os.path.dirname(__file__), 'build')
    os.makedirs(build_dir, exist_ok=True)

    llm_txt_path = os.path.join(build_dir, 'llms.txt')
    with open(llm_txt_path, 'w', encoding='utf-8') as f:
        f.write(content)

    print(f"Generated llm.txt with {len(docs_files)} docs and {len(optional_files)} optional files")


def on_post_build(config):
    # You can pass custom text here
    custom_description = "This is the comprehensive documentation for MACH composer, a framework for orchestrating modern cloud-native commerce architectures."
    generate_llm_txt(nav=config.nav, custom_text=custom_description)
