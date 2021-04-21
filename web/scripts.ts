import positions from "../positions.json";

// elem = element to create
// html = tag's html contents
// tags = [["attribute", "value"], ...]
function createEl(elem: string, html: string = null, tags = []): HTMLElement {
  const el = document.createElement(elem);
  el.innerHTML = html;
  tags.forEach((tag) => {
    if (
      tag[0] === undefined ||
      tag[0] === null ||
      tag[1] === undefined ||
      tag[1] === null
    ) {
      return;
    }

    el.setAttribute(tag[0], tag[1]);
  });
  return el;
}

function createMetadata(string) {
  const parts = string.split("_");

  if (parts.length != 3) {
    return { name: "", position: "", suppl: "" };
  }

  return {
    name: parts[0],
    position: parts[1],
    suppl: parts[2].replaceAll("-", "/").replace(/\.[^/.]+$/, ""), // trims extension,
  };
}

function createInput(value, name, label) {
  const div = createEl("div");
  const labelEl = createEl("label", label + ": ");
  labelEl.setAttribute("for", name);
  const input = createEl("input", null, [
    ["type", "text"],
    ["value", value],
    ["name", name],
    ["id", name],
  ]);
  div.appendChild(labelEl);
  div.appendChild(input);
  return div;
}

function createPositionsDropdown(value, name) {
  const div = createEl("div");
  const labelEl = createEl("label", "Stilling: ");
  labelEl.setAttribute("for", name);
  const input = createEl("select", null, [
    ["name", name],
    ["id", name],
  ]);

  // Loop over positions
  positions.forEach((pos) => {
    input.appendChild(
      createEl("option", pos.title, [
        ["value", pos.abbr],
        ["selected", pos.abbr === value ? "true" : null],
      ])
    );
  });

  div.appendChild(labelEl);
  div.appendChild(input);
  return div;
}

export function load(el: HTMLInputElement) {
  document.querySelector("#images").innerHTML = null;

  const files = el.files;
  const rows = createEl("div");
  for (let n = 0; n < files.length; n++) {
    const row = document.createElement("fieldset");
    const person = createMetadata(files[n].name);

    row.appendChild(createEl("legend", files[n].name));

    const reader = new FileReader();
    reader.onload = function () {
      row.appendChild(
        createEl("img", null, [
          ["src", reader.result],
          ["width", 100],
        ])
      );
    };
    reader.readAsDataURL(files[n]);

    const div = createEl("div", null, [["class", "person-inputs"]]);
    div.appendChild(createInput(person.name, files[n].name + "-name", "Navn"));
    div.appendChild(
      createPositionsDropdown(person.position, files[n].name + "-position")
    );
    div.appendChild(createInput(person.suppl, files[n].name + "-suppl", "Supplerende"));
    row.appendChild(div);

    rows.appendChild(row);
  }

  document.querySelector("#images").appendChild(rows);
}

declare global {
  interface Window {
    load: (el: HTMLInputElement) => void;
  }
}

(() => {
  window.load = load;
})();
