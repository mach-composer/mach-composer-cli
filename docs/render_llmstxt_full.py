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


def generate_llm_txt(nav):
    """Generate llm.txt file from the navigation structure"""
    import os

    # Extract all files from nav
    all_files = extract_files_from_nav(nav)

    # Start building the full content
    full_content = ""

    build_dir = os.path.join(os.path.dirname(__file__), 'build')

    # Iterate through all files and read their content
    for file_path in all_files:
        src_file_path = os.path.join(build_dir, file_path)

        # Check if file exists in src directory
        if not os.path.exists(src_file_path):
            raise Exception(f"File {src_file_path} does not exist")

        with open(src_file_path, 'r', encoding='utf-8') as f:
            file_content = f.read()

        # Split content into lines
        lines = file_content.split('\n')

        if lines and lines[0].startswith('# '):
            # Insert Source line after the title
            modified_content = lines[0] + '\n' + f'Source: {file_path}\n' + '\n'.join(lines[1:])
            full_content += modified_content + '\n\n'
        else:
            # Fallback if no title found - add source at the beginning
            full_content += f'Source: {file_path}\n'
            full_content += file_content + '\n\n'

    # Write to build folder
    os.makedirs(build_dir, exist_ok=True)

    llm_full_txt_path = os.path.join(build_dir, 'llms-full.txt')
    with open(llm_full_txt_path, 'w', encoding='utf-8') as f:
        f.write(full_content)

    print(f"Generated llm-full.txt with content from {len(all_files)} files")


def on_post_build(config):
    generate_llm_txt(nav=config.nav)
