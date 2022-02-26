Overview
--------

This repository contains the BibTeX file and HTML templates that form the
[Internet Censorship Bibliography](https://censorbib.nymity.ch).

Build it
--------

You first need [`bibliogra.py`](https://github.com/NullHypothesis/bibliograpy)
to turn the BibTeX file into an HTML bibliography.

Then, run the following commands to write the bibliography to `OUTPUT_DIR`.

    $ ./fetch_pdfs.py references.bib OUTPUT_DIR
    $ bibliogra.py -H header.tpl -F footer.tpl -f references.bib OUTPUT_DIR

Acknowledgements
----------------

CensorBib uses [Font Awesome](https://fontawesome.com/license/free) icons
without modification.

Feedback
--------

Contact: Philipp Winter <phw@nymity.ch>
