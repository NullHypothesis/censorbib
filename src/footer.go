package main

func footer() string {
	return `<div id="bibtex-modal" hidden>
  <div id="bibtex-backdrop" data-close-bibtex></div>
  <section id="bibtex-dialog" role="dialog" aria-modal="true" aria-labelledby="bibtex-title">
    <header>
      <h2 id="bibtex-title">BibTeX</h2>
      <span id="bibtex-copy-status" aria-live="polite"></span>
      <button id="bibtex-copy" type="button">Copy</button>
      <button id="bibtex-close" type="button" data-close-bibtex>Close</button>
    </header>
    <pre id="bibtex-content"></pre>
  </section>
</div>

<div id="footer">
Icons taken without modification from
<a href="https://fontawesome.com/license">Font Awesome</a>.
</div>

<script>
(function () {
  const dataElement = document.getElementById("reference-data");
  const references = dataElement ? JSON.parse(dataElement.textContent) : [];
  const referencesByCiteName = new Map(references.map((reference) => [reference.citeName, reference]));
  const form = document.getElementById("search-form");
  const input = document.getElementById("search-input");
  const resultCount = document.getElementById("result-count");
  const noResults = document.getElementById("no-results");
  const modal = document.getElementById("bibtex-modal");
  const modalTitle = document.getElementById("bibtex-title");
  const modalContent = document.getElementById("bibtex-content");
  const copyButton = document.getElementById("bibtex-copy");
  const copyStatus = document.getElementById("bibtex-copy-status");
  const closeButton = document.getElementById("bibtex-close");

  function normalize(value) {
    return String(value || "")
      .toLocaleLowerCase("en")
      .normalize("NFD")
      .replace(/[\u0300-\u036f]/g, "")
      .replace(/[^\p{L}\p{N}]+/gu, " ")
      .trim();
  }

  function tokenize(value) {
    const normalized = normalize(value);
    return normalized === "" ? [] : normalized.split(/\s+/);
  }

  for (const reference of references) {
    reference.searchTokens = tokenize([
      reference.citeName,
      reference.title,
      reference.authors,
      reference.venue,
      reference.year,
      reference.publisher,
    ].join(" "));
    reference.item = document.getElementById(reference.citeName);
    reference.originalHTML = reference.item ? reference.item.innerHTML : "";
  }

  function referenceMatches(reference, queryTokens) {
    return queryTokens.every((queryToken) =>
      reference.searchTokens.some((token) => token.startsWith(queryToken) || token.includes(queryToken))
    );
  }

  function pluralize(count, singular, plural) {
    return count === 1 ? singular : plural;
  }

  function updateURL(query) {
    const params = new URLSearchParams(window.location.search);
    if (query === "") {
      params.delete("q");
    } else {
      params.set("q", query);
    }
    const search = params.toString();
    const nextURL = window.location.pathname + (search ? "?" + search : "") + window.location.hash;
    window.history.replaceState({}, "", nextURL);
  }

  function tokenizeForHighlight(value) {
    return Array.from(new Set(
      String(value || "")
        .toLocaleLowerCase("en")
        .split(/[^\p{L}\p{N}]+/gu)
        .filter(Boolean)
    )).sort((left, right) => right.length - left.length);
  }

  function nextMatch(text, lowerText, start, highlightTokens) {
    let matchIndex = -1;
    let matchToken = "";

    for (const token of highlightTokens) {
      const index = lowerText.indexOf(token, start);
      if (index === -1) {
        continue;
      }
      if (matchIndex === -1 || index < matchIndex || (index === matchIndex && token.length > matchToken.length)) {
        matchIndex = index;
        matchToken = token;
      }
    }

    return { index: matchIndex, token: matchToken };
  }

  function highlightTextNode(textNode, highlightTokens) {
    const text = textNode.nodeValue;
    const lowerText = text.toLocaleLowerCase("en");
    const fragment = document.createDocumentFragment();
    let index = 0;

    while (index < text.length) {
      const match = nextMatch(text, lowerText, index, highlightTokens);
      if (match.index === -1) {
        fragment.append(document.createTextNode(text.slice(index)));
        break;
      }

      if (match.index > index) {
        fragment.append(document.createTextNode(text.slice(index, match.index)));
      }

      const mark = document.createElement("mark");
      mark.textContent = text.slice(match.index, match.index + match.token.length);
      fragment.append(mark);
      index = match.index + match.token.length;
    }

    textNode.replaceWith(fragment);
  }

  function highlightMatches(item, highlightTokens) {
    if (highlightTokens.length === 0) {
      return;
    }

    const walker = document.createTreeWalker(item, NodeFilter.SHOW_TEXT, {
      acceptNode(node) {
        if (!node.nodeValue.trim() || node.parentElement.closest(".icons")) {
          return NodeFilter.FILTER_REJECT;
        }
        return NodeFilter.FILTER_ACCEPT;
      },
    });
    const textNodes = [];
    while (walker.nextNode()) {
      textNodes.push(walker.currentNode);
    }
    for (const textNode of textNodes) {
      highlightTextNode(textNode, highlightTokens);
    }
  }

  function applySearch(query, shouldUpdateURL) {
    const queryTokens = tokenize(query);
    const highlightTokens = tokenizeForHighlight(query);
    let visibleCount = 0;

    for (const reference of references) {
      const item = reference.item;
      if (!item) {
        continue;
      }
      item.innerHTML = reference.originalHTML;
      const visible = queryTokens.length === 0 || referenceMatches(reference, queryTokens);
      item.hidden = !visible;
      if (visible) {
        highlightMatches(item, highlightTokens);
        visibleCount++;
      }
    }

    for (const group of document.querySelectorAll(".year-group")) {
      group.hidden = !group.querySelector("li:not([hidden])");
    }

    resultCount.textContent = visibleCount + " " + pluralize(visibleCount, "paper", "papers");
    noResults.hidden = visibleCount !== 0;
    if (shouldUpdateURL) {
      updateURL(query.trim());
    }
  }

  input.addEventListener("input", () => {
    applySearch(input.value, true);
  });

  form.addEventListener("submit", (event) => {
    event.preventDefault();
    applySearch(input.value, true);
  });

  function openBibtex(citeName) {
    const reference = referencesByCiteName.get(citeName);
    if (!reference) {
      return;
    }
    modalTitle.textContent = citeName;
    modalContent.textContent = reference.rawBibtex;
    copyStatus.textContent = "";
    modal.hidden = false;
    copyButton.focus();
  }

  function closeBibtex() {
    modal.hidden = true;
  }

  document.addEventListener("click", (event) => {
    const bibtexLink = event.target.closest(".bibtex-link");
    if (bibtexLink) {
      event.preventDefault();
      openBibtex(bibtexLink.dataset.reference);
      return;
    }
    if (event.target.closest("[data-close-bibtex]")) {
      closeBibtex();
    }
  });

  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && !modal.hidden) {
      closeBibtex();
    }
  });

  copyButton.addEventListener("click", async () => {
    try {
      await navigator.clipboard.writeText(modalContent.textContent);
      copyStatus.textContent = "Copied";
    } catch (error) {
      copyStatus.textContent = "Copy failed";
    }
  });

  closeButton.addEventListener("click", closeBibtex);

  const initialQuery = new URLSearchParams(window.location.search).get("q") || "";
  input.value = initialQuery;
  applySearch(initialQuery, false);
})();
</script>

</body>
</html>`
}
