import os


def on_page_markdown(markdown, page, **kwargs):
    # Get the directory where this render_markdown.py script is located
    script_dir = os.path.dirname(os.path.abspath(__file__))
    url = script_dir + "/build/" + page.url
    if url:
        # Replace .html with .md to get the output file path
        if url.endswith(".html"):
            output_path = url.replace(".html", ".md")
        else:
            # If no .html extension, just append .md
            output_path = url + ".md" if not url.endswith(".md") else url

        # Create directory if it doesn't exist
        os.makedirs(os.path.dirname(output_path), exist_ok=True)

        # Write markdown content to the file
        with open(output_path, "w", encoding="utf-8") as f:
            f.write(markdown)
