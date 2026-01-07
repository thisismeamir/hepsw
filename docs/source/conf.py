# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#project-information

project = 'hepsw'
copyright = '2025, Amir H. Ebrahimnezhad'
author = 'Amir H. Ebrahimnezhad'
release = '0.0.5'

extensions = [
    "myst_parser",
    "sphinx.ext.graphviz"
]

html_theme = "alabaster"
graphviz_output_format = "svg"
source_suffix = {
    ".rst": "restructuredtext",
    ".md": "markdown",
}

myst_enable_extensions = [
    "colon_fence",
    "deflist",
]

