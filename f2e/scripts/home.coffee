$ = require 'jquery'

window.jQuery = $  # Orz....... for trouble of bootstrap

require 'bootstrap'
require 'bootstrap-select'
require 'basil'

$ ->
  $('.j-select').selectpicker()

$ ->
  basil = window.Basil()

  submitButton = $('.j-submit')
  submitButton.data('value', submitButton.val())

  counter = basil.get('counter') || 0
  counting = () ->
    if counter >= 0
      submitButton.val("#{submitButton.data('value')} (#{counter})")
      submitButton.addClass('disabled')
    else
      submitButton.val(submitButton.data('value'))
      submitButton.removeClass('disabled')
    counter--
    basil.set('counter', counter)
  start = () ->
    counting()
    setInterval(counting, 1000)

  if counter > 0
    start()
  if $('.j-flashed-message[data-category="success"]').length > 0
    counter = 60
    basil.set('counter', counter)
    start()
