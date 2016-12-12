function check_task_status(task_id, trial, max_trial, timeout) {
  if (trial >= max_trial) {
    $('#response').attr('class', 'fail').html('Задача не завершилась');
    return;
  }
  $.getJSON('/digest/check_task_status', {'task_id':task_id})
    .done(function(data) {
      if (data['success'] === false) {
        if (data['msg'] === 'pending') {
          setTimeout(check_task_status.bind(null, task_id, ++trial, max_trial, timeout), timeout);
        }
        else {
          $('#response').attr('class', 'fail').html(data['msg']);
        }
      }
      else {
        $('#response').attr('class', 'success').html(data['msg']);
      }
    })
    .fail(function(data) {
      $('#response').attr('class', 'fail').html('Ошибка ajax');
    });
}


$(document).ready(function() {
  $('#send_digest_to_email').submit(function(event) {
    event.preventDefault();
    var pResponse = $('#response');
    pResponse.removeAttr('class').html('Ждём результата...');
    $.post($(this).attr('action'), $(this).serialize())
      .done(function(data) {
        data = JSON.parse(data);
        if (data['success'] === true) {
          var task_id = data['msg'];
          var timeout = 3000;
          setTimeout(check_task_status.bind(null, task_id, 1, 10, timeout), timeout);
        }
        else {
          pResponse.attr('class', 'fail');
        }
      })
      .fail(function() {
        pResponse.attr('class', 'fail').html('Ошибка ajax!');
      });
  });
})
