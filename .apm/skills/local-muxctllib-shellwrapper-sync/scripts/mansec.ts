#!/usr/bin/env -S deno run --allow-read
/**
 * Print one command section of a rendered tmux man page as Go `//` comment lines.
 *
 * usage: deno run --allow-read mansec.ts <tmux.man> <command>
 *
 * Render the page first with hyphenation and justification off, so the copy is verbatim:
 *
 *     MANWIDTH=90 MANPAGER=cat man --nh --nj tmux | sed 's/.\x08//g' > tmux.man
 */

const [path, cmd] = Deno.args;
if (!path || !cmd) {
  console.error("usage: mansec.ts <tmux.man> <command>");
  Deno.exit(2);
}

const lines = Deno.readTextFileSync(path).split("\n");
const escaped = cmd.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
const head = new RegExp(`^ {5}${escaped}(\\s|$)`);

// A command name also appears at line start in prose and in the key binding table, so pick
// the synopsis: the one followed by "(alias: ...)", else the one listing flags in brackets.
const cands = lines.flatMap((l, i) => (head.test(l) ? [i] : []));
const start = cands.find((i) => lines.slice(i, i + 6).some((x) => x.includes("(alias:"))) ??
  cands.find((i) => lines[i].includes("["));
if (start === undefined) {
  console.error(`not found: ${cmd}`);
  Deno.exit(1);
}

const out = [lines[start]];
for (const l of lines.slice(start + 1)) {
  // the next command synopsis (5-space indent) or a section header (column 0) ends it
  if (/^ {5}\S/.test(l) || /^\S/.test(l)) break;
  out.push(l);
}
while (out.length > 0 && out[out.length - 1].trim() === "") out.pop();

for (const l of out) {
  const body = l.startsWith("     ") ? l.slice(5) : l.trimStart();
  console.log(`// ${body}`.trimEnd());
}
