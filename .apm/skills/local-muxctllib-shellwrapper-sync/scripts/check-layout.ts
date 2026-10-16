#!/usr/bin/env -S deno run --allow-read
/**
 * Check the mandatory layout of shellwrapper command files.
 *
 * usage: deno run --allow-read check-layout.ts <shellwrapper dir>
 *
 * For every command file (not shell.go, not tests) checks:
 *   - a free-floating comment follows the imports and is followed by a blank line;
 *   - sections come in order: option types, result types, the exported ShellOut method, then
 *     helpers (unexported funcs and methods, including unexported ShellOut methods);
 *   - there is one exported ShellOut method, and its option parameter is <Method>Option.
 * Exits 1 if anything is reported.
 */

const root = Deno.args[0];
if (!root) {
  console.error("usage: check-layout.ts <shellwrapper dir>");
  Deno.exit(2);
}

const OPTION = 1, RESULT = 2, METHOD = 3, HELPER = 4;
let bad = 0;

function report(file: string, msg: string) {
  bad++;
  console.log(`${file}: ${msg}`);
}

function checkFile(file: string) {
  const L = Deno.readTextFileSync(file).split("\n");

  let i = L.findIndex((l) => l.startsWith("import"));
  if (L[i].trimEnd().endsWith("(")) {
    while (L[i] !== ")") i++;
  }
  i++;
  while (L[i] === "") i++;
  if (!L[i].startsWith("//")) {
    report(file, "no free-floating man/help comment after imports");
    return;
  }
  let end = i;
  while (L[end].startsWith("//")) end++;
  if (L[end] !== "") report(file, "man/help comment is not followed by a blank line");

  // A type embedded by an option struct (e.g. OptionScope) is itself an option type.
  const embedded = new Set(L.flatMap((l) => l.match(/^\t([A-Z]\w*)$/)?.[1] ?? []));
  const kinds: [number, string][] = [];
  for (const l of L.slice(end)) {
    const t = l.match(/^type (\w+) struct/);
    if (t && (t[1].endsWith("Option") || embedded.has(t[1]))) kinds.push([OPTION, l]);
    else if (l.startsWith("type ")) kinds.push([RESULT, l]);
    else if (/^func \(s \*ShellOut\) [A-Z]/.test(l)) kinds.push([METHOD, l]);
    else if (l.startsWith("func ")) kinds.push([HELPER, l]);
  }
  if (kinds.some(([k], n) => n > 0 && k < kinds[n - 1][0])) {
    report(file, "section order: " + kinds.map(([, l]) => l.slice(0, 50)).join("; "));
  }

  const methods = kinds.filter(([k]) => k === METHOD).map(([, l]) => l);
  if (methods.length !== 1) {
    report(file, `want 1 exported ShellOut method, got ${methods.length}`);
    return;
  }
  const m = methods[0].match(/^func \(s \*ShellOut\) (\w+)\(ctx context\.Context(, opt (\w+))?/);
  if (m?.[3] && m[3] !== `${m[1]}Option`) {
    report(file, `option type ${m[3]}, want ${m[1]}Option`);
  }
}

const files: string[] = [];
for (const pkg of Deno.readDirSync(root)) {
  if (!pkg.isDirectory) continue;
  for (const e of Deno.readDirSync(`${root}/${pkg.name}`)) {
    if (!e.isFile || !e.name.endsWith(".go")) continue;
    if (e.name.endsWith("_test.go") || e.name === "shell.go") continue;
    files.push(`${root}/${pkg.name}/${e.name}`);
  }
}
for (const f of files.sort()) checkFile(f);

Deno.exit(bad > 0 ? 1 : 0);
