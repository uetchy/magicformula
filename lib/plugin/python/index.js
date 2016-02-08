'use strict';

var fs = require('fs');
var ejs = require('ejs');

function parseDeps(content) {
  return [];
}

module.exports = {
  props: [{
    name: 'target-path-64',
    message: 'Target binary file with 64bit arch',
    required: true
  }],
  call: function call(props) {
    var baseDir = props.baseDir;
    var deps = fs.readFileSync(path.join(baseDir, 'requirements.txt'));

    deps.forEach(function (dep) {
      formula.resource(dep.name, dep.url, sha256(dep.path));
    });

    var template = fs.readFileSync('./src/plugin/binary/template.ejs', 'utf-8');
    return ejs.render(template, { props: props });
  }
};