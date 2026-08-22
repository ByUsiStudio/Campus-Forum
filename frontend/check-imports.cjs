const fs = require('fs')
const path = require('path')

const missing = []
function walk(d) {
  for (const e of fs.readdirSync(d, { withFileTypes: true })) {
    const p = path.join(d, e.name)
    if (e.isDirectory()) {
      if (e.name !== 'node_modules') walk(p)
    } else if (/\.(vue|js)$/.test(e.name)) {
      const dir = path.dirname(p)
      const t = fs.readFileSync(p, 'utf8')
      const re = /(?:from\s+['"]|import\s+['"]|require\(['"])([^'"]+)['"]/g
      let m
      while ((m = re.exec(t))) {
        const spec = m[1]
        if (spec.startsWith('@/')) {
          let rel = path.join('src', spec.slice(2))
          if (!resolves(rel)) missing.push(p.toString() + ' -> ' + spec)
        } else if (spec.startsWith('./') || spec.startsWith('../')) {
          let rel = path.resolve(dir, spec)
          if (!resolves(rel) && !/\.(css|png|svg|jpg)$/.test(spec)) missing.push(p.toString() + ' -> ' + spec)
        }
      }
    }
  }
}
function resolves(rel) {
  try {
    if (/\.(vue|js|css|mjs|json)$/.test(rel) || fs.existsSync(rel)) {
      return fs.existsSync(rel)
    }
    return fs.existsSync(rel + '.js') || fs.existsSync(rel + '.vue') || fs.existsSync(rel + '/index.js')
  } catch (e) { return true }
}
walk('src')
if (missing.length) {
  console.log('MISSING IMPORTS (' + missing.length + '):')
  missing.forEach(x => console.log('  ' + x))
} else {
  console.log('All internal imports resolve. OK')
}
