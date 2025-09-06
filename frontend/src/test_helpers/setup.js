// This file is included by mocha to help set up testing

// First set up sinon-chai
const sinonChai = require('sinon-chai')
const chai = require('chai')
chai.use(sinonChai.default)

const Module = require('module')
const path = require('path')

// Hook into NodeJS module resolution to stub out webpack-handled imports such as
// global src import and importing of assets. This way tests won't fail when testing
// a file that imports a component
function resolveImport(request) {
  if (request.startsWith('src/')) return path.join(__dirname, '../..', request)
  if (request === './stylesheet') return request = path.join(__dirname, 'mock-stylesheet.js')
  return request
}
const defaultFilenameResolver = Module._resolveFilename
require('module')._resolveFilename = function(request, parentModule, isMain, options) {
  return defaultFilenameResolver(resolveImport(request), parentModule, isMain, options)
}
