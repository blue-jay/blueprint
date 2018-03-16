'use strict'

var PassThrough = require('readable-stream').PassThrough
var inherits = require('inherits')
var p = require('process-nextick-args')

function Cloneable (stream, opts) {
  if (!(this instanceof Cloneable)) {
    return new Cloneable(stream, opts)
  }

  var objectMode = stream._readableState.objectMode
  this._original = stream
  this._clonesCount = 1

  opts = opts || {}
  opts.objectMode = objectMode

  PassThrough.call(this, opts)

  forwardDestroy(stream, this)

  this.on('newListener', onData)
  this.once('resume', onResume)
}

inherits(Cloneable, PassThrough)

function onData (event, listener) {
  if (event === 'data' || event === 'readable') {
    this.removeListener('newListener', onData)
    this.removeListener('resume', onResume)
    p.nextTick(clonePiped, this)
  }
}

function onResume () {
  this.removeListener('newListener', onData)
  p.nextTick(clonePiped, this)
}

Cloneable.prototype.clone = function () {
  if (!this._original) {
    throw new Error('already started')
  }

  this._clonesCount++

  // the events added by the clone should not count
  // for starting the flow
  this.removeListener('newListener', onData)
  var clone = new Clone(this)
  this.on('newListener', onData)

  return clone
}

function forwardDestroy (src, dest) {
  src.on('error', destroy)
  src.on('close', destroy)

  function destroy (err) {
    dest.destroy(err)
  }
}

function clonePiped (that) {
  if (--that._clonesCount === 0 && !that._readableState.destroyed) {
    that._original.pipe(that)
    that._original = undefined
  }
}

function Clone (parent, opts) {
  if (!(this instanceof Clone)) {
    return new Clone(parent, opts)
  }

  var objectMode = parent._readableState.objectMode

  opts = opts || {}
  opts.objectMode = objectMode

  this.parent = parent

  PassThrough.call(this, opts)

  forwardDestroy(this.parent, this)

  parent.pipe(this)

  // the events added by the clone should not count
  // for starting the flow
  // so we add the newListener handle after we are done
  this.on('newListener', onDataClone)
  this.on('resume', onResumeClone)
}

function onDataClone (event, listener) {
  // We start the flow once all clones are piped or destroyed
  if (event === 'data' || event === 'readable' || event === 'close') {
    p.nextTick(clonePiped, this.parent)
    this.removeListener('newListener', onDataClone)
  }
}

function onResumeClone () {
  this.removeListener('newListener', onDataClone)
  p.nextTick(clonePiped, this.parent)
}

inherits(Clone, PassThrough)

Clone.prototype.clone = function () {
  return this.parent.clone()
}

Cloneable.isCloneable = function (stream) {
  return stream instanceof Cloneable || stream instanceof Clone
}

module.exports = Cloneable
