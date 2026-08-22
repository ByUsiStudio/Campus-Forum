/* Static SFC verifier using @vue/compiler-sfc (pure JS, no esbuild spawn) */
import { readFileSync, readdirSync, statSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import Path from 'node:path'
import { parse, compileScript, compileTemplate } from '@vue/compiler-sfc'

const src = Path.resolve(Path.dirname(fileURLToPath(import.meta.url)), '../src')
const files = []
function walk(dir) {
  for (const e of readdirSync(dir)) {
    const p = Path.join(dir, e)
    const st = statSync(p)
    if (st.isDirectory()) walk(p)
    else if (e.endsWith('.vue')) files.push(p)
  }
}
walk(src)
console.log(`Found ${files.length} .vue files`)

let errors = 0
let parsed = 0
for (const f of files) {
  const source = readFileSync(f, 'utf8')
  const { descriptor } = parse(source, { filename: f })
  let ok = true
  if (descriptor.errors && descriptor.errors.length) {
    for (const e of descriptor.errors) { console.log(`[PARSE] ${Path.relative(src, f)}: ${e.message}`); errors++ }
    ok = false
  }
  if (ok) {
    const id = Buffer.from(f).toString('base64')
    try {
      if (descriptor.script || descriptor.scriptSetup) {
        compileScript(descriptor, { id })
      }
    } catch (e) {
      console.log(`[SCRIPT] ${Path.relative(src, f)}: ${e.message}`)
      errors++
    }
    if (descriptor.template && descriptor.template.content.trim()) {
      try {
        const res = compileTemplate({
          source: descriptor.template.content,
          filename: f,
          id,
          compilerOptions: { bindingMetadata: {}, mode: 'module' }
        })
        if (res.errors && res.errors.length) {
          for (const e of res.errors) { console.log(`[TEMPLATE] ${Path.relative(src, f)}: ${e.message}`); errors++ }
        }
      } catch (e) {
        console.log(`[TEMPLATE-CRASH] ${Path.relative(src, f)}: ${e.message}`)
        errors++
      }
    }
  }
  parsed++
}

console.log(`\nDone. Parsed ${parsed} files, ${errors} error(s).`)
process.exit(errors ? 1 : 0)
