import '@testing-library/jest-dom'

// JSDOM has no dialog top layer. Model only the open/close API here;
// focus containment and background inertness need real-browser verification.
HTMLDialogElement.prototype.showModal = function () {
  this.setAttribute('open', '')
}
HTMLDialogElement.prototype.close = function () {
  this.removeAttribute('open')
}
