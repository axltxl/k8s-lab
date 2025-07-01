#!/usr/bin/env python3

# File: scripts/remote/j2.py
# Description: A script to render Jinja2 templates using variables from a YAML file.

import sys
import yaml
from jinja2 import Environment, FileSystemLoader

def main(yaml_file, template_file):
    # Load variables from YAML
    with open(yaml_file, 'r') as f:
        variables = yaml.safe_load(f)

    # Load template from the same directory as the template file
    template_dir = template_file.rsplit('/', 1)[0] or '.'
    template_name = template_file.rsplit('/', 1)[-1]

    env = Environment(loader=FileSystemLoader(template_dir))
    template = env.get_template(template_name)

    # Render template with variables
    output = template.render(variables)
    print(output)

if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: j2 <variables.yaml> <template.j2>")
        sys.exit(1)

    yaml_file = sys.argv[1]
    template_file = sys.argv[2]
    main(yaml_file, template_file)
