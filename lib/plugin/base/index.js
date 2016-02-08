'use strict';

module.exports = {
  props: [{
    name: 'name',
    pattern: /^[a-zA-Z\s\-]+$/,
    message: 'Package name',
    required: true
  }, {
    name: 'owner',
    type: 'string',
    message: 'Owner name',
    required: true
  }, {
    name: 'tag',
    type: 'string',
    message: 'Product tag/version'
  }],
  call: function call(props, formula) {
    formula.name = props.name;
    formula.owner = props.owner;
    formula.tag = props.tag;
    return formula;
  }
};