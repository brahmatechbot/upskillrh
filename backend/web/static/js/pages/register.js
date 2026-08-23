$(function () {
  const $root = $('.wide-card');
  const initialType = $root.data('initial-type');
  const $form = $('#register-form');
  const $summary = $('#form-summary');
  const $submit = $('#submit-button');
  let submitting = false;

  function accountType() { return $('input[name="account_type"]:checked').val() || ''; }
  function showSummary(message) { $summary.prop('hidden', !message).text(message || ''); if (message) $summary.trigger('focus'); }
  function clearErrors() { $('.field-error').text(''); }
  function setErrors(fields) { clearErrors(); Object.keys(fields || {}).forEach((k) => $('#' + k + '-error').text(fields[k])); }
  function updateType() {
    const type = accountType();
    $('.candidate-only').toggle(type === 'candidate');
    $('.company-only, .company-label').toggle(type === 'company');
    $('#other-segment-wrap').toggle(type === 'company' && $('#industry_segment_code').val() === 'other');
  }

  if (initialType === 'empresa') $('input[value="company"]').prop('checked', true);
  if (initialType === 'candidato') $('input[value="candidate"]').prop('checked', true);
  updateType();

  $('input[name="account_type"]').on('change', function () {
    const previous = $(this).val() === 'candidate' ? 'company' : 'candidate';
    if (previous === 'company') $('.company-only input, .company-only select').val('');
    if (previous === 'candidate') $('.candidate-only input').val('');
    updateType();
  });
  $('#industry_segment_code').on('change', updateType);
  $('#toggle-password').on('click', function () { const $p = $('#password'); const hidden = $p.attr('type') === 'password'; $p.attr('type', hidden ? 'text' : 'password'); $(this).text(hidden ? 'Ocultar' : 'Mostrar'); });

  $form.on('submit', function (event) {
    event.preventDefault();
    if (submitting) return;
    const payload = {
      account_type: accountType(),
      full_name: $('#full_name').val().trim(),
      responsible_name: $('#responsible_name').val().trim(),
      email: $('#email').val().trim(),
      password: $('#password').val(),
      password_confirmation: $('#password_confirmation').val(),
      cpf: $('#cpf').val().trim(),
      trade_name: $('#trade_name').val().trim(),
      employee_range_code: $('#employee_range_code').val(),
      cnpj: $('#cnpj').val().trim(),
      industry_segment_code: $('#industry_segment_code').val(),
      other_industry_segment: $('#other_industry_segment').val().trim(),
      accept_terms: $('#accept_terms').is(':checked'),
      accept_privacy: $('#accept_privacy').is(':checked'),
      marketing_granted: $('#marketing_granted').is(':checked')
    };
    submitting = true; showSummary(''); clearErrors(); $submit.prop('disabled', true).addClass('is-loading').find('.button-text').text('Criando...');
    $.ajax({ url: '/api/v1/auth/register', method: 'POST', contentType: 'application/json', headers: { 'X-CSRF-Token': $('#csrf_token').val() }, data: JSON.stringify(payload) })
      .done((response) => window.location.assign(response.data.next_url))
      .fail((xhr) => { const body = xhr.responseJSON || {}; const error = body.error || {}; setErrors(error.fields || {}); showSummary(error.message || 'Não foi possível criar a conta agora.'); })
      .always(() => { submitting = false; $submit.prop('disabled', false).removeClass('is-loading').find('.button-text').text('Criar conta'); });
  });
});
