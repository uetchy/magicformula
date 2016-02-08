function calcSHA256(src) {
  return "sha256-hash";
}

module.exports = {
  props: [{
    name: 'target-path',
    message: 'target binary file',
    required: true
  }],
  call: function(props, formula) {
    formula.url = "https://github.com/" + props.owner + "/" + props.name + "/releases/download/" + props.tag;
    formula.sha256 = calcSHA256(props.targetPath);
    return formula;
  }
};
