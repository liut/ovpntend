$ = require 'jquery'

window.jQuery = $  # Orz....... for trouble of bootstrap

require 'bootstrap'
require 'bootstrap-select'
require 'basil.js/build/basil.js'
require 'bootstrap-validator'

$ ->
  $('.j-select').selectpicker()

$ ->
  basil = new window.Basil()

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

$ ->
  $('#form1').validator().on 'submit', (e) ->
    if e.isDefaultPrevented()
      # handle the invalid form...
      console.log(e)
    else
      # everything looks good!
      # console.log(e)
      data = {
        name: $('#emailName').val(),
        host: $('#emailHost').val(),
        oscat: $('input[name=oscat]:checked').val()
      }
      $.post e.target.action, data, (res) ->
        console.log(res)
        alert('ok! please checke your email.')

    return false
