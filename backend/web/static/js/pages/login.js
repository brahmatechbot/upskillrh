$(function () {
  const $form = $('#login-form');
  const $email = $('#email');
  const $password = $('#password');
  const $summary = $('#form-summary');
  const $submit = $('#submit-button');
  let submitting = false;

  function setError(field, message) {
    $('#' + field + '-error').text(message || '');
  }

  function showSummary(message) {
    if (!message) {
      $summary.prop('hidden', true).text('');
      return;
    }
    $summary.prop('hidden', false).text(message).trigger('focus');
  }

  function validate() {
    const email = $email.val().trim();
    const password = $password.val();
    const errors = [];
    setError('email', '');
    setError('password', '');

    if (!email) {
      setError('email', 'Informe seu e-mail.');
      errors.push('Informe seu e-mail.');
    } else if (email.length > 254 || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      setError('email', 'Informe um e-mail válido.');
      errors.push('Informe um e-mail válido.');
    }

    if (!password) {
      setError('password', 'Informe sua senha.');
      errors.push('Informe sua senha.');
    }

    if (errors.length) {
      showSummary(errors.join(' '));
      return false;
    }
    showSummary('');
    return true;
  }

  $('#toggle-password').on('click', function () {
    const isHidden = $password.attr('type') === 'password';
    $password.attr('type', isHidden ? 'text' : 'password');
    $(this).attr('aria-pressed', String(isHidden)).text(isHidden ? 'Ocultar' : 'Mostrar');
    $password.trigger('focus');
  });

  $form.on('submit', function (event) {
    event.preventDefault();
    if (submitting || !validate()) return;

    submitting = true;
    $submit.prop('disabled', true).addClass('is-loading').find('.button-text').text('Entrando...');

    $.ajax({
      url: '/api/v1/auth/login',
      method: 'POST',
      contentType: 'application/json',
      headers: { 'X-CSRF-Token': $('#csrf_token').val() },
      data: JSON.stringify({
        email: $email.val().trim(),
        password: $password.val(),
        remember_me: $('#remember_me').is(':checked')
      })
    }).done(function (response) {
      window.location.assign(response.data.next_url);
    }).fail(function (xhr) {
      const body = xhr.responseJSON || {};
      const error = body.error || {};
      const fields = error.fields || {};
      setError('email', fields.email || '');
      setError('password', fields.password || '');
      let message = error.message || 'Não foi possível entrar agora. Tente novamente.';
      if (error.code === 'INTERNAL_ERROR' && body.meta && body.meta.request_id) {
        message += ' Código: ' + body.meta.request_id;
      }
      showSummary(message);
    }).always(function () {
      submitting = false;
      $submit.prop('disabled', false).removeClass('is-loading').find('.button-text').text('Entrar');
    });
  });
});
