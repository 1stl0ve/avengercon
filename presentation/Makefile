
DATE_COVER="December 1, 2020"
SOURCE_FORMAT="markdown_github+yaml_metadata_block+implicit_figures+all_symbols_escapable+link_attributes+smart+fenced_divs"
DATA_DIR="pandoc"

slides:
	pandoc \
		--verbose \
		--dpi=300 \
		--slide-level 2 \
		--shift-heading-level=0 \
		--data-dir="${DATA_DIR}" \
		--template default_mod.latex \
		--pdf-engine latexmk \
		-f ${SOURCE_FORMAT} \
		-M date=${DATE_COVER} \
		-V classoption:aspectratio=169 \
		-V lang=en-US \
		-t beamer \
		slides.md \
		-o slides.pdf

# pandoc -s --dpi=300 --slide-level 2 --toc --listings --shift-heading-level=0 --columns=50 --template default_mod.latex --pdf-engine xelatex -f "$SOURCE_FORMAT" -M date="$DATE_COVER" -V classoption:aspectratio=169 -V lang=en-US -t beamer presentation.md -o presentation.pdf

.PHONY: slides