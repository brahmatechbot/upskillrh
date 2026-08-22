$(function () {
  $('.reveal').addClass('visible');

  $('a[href^="#"]').on('click', function (event) {
    const targetSelector = $(this).attr('href');
    if (!targetSelector || targetSelector === '#') return;

    const $target = $(targetSelector);
    if (!$target.length) return;

    event.preventDefault();
    $('html, body').animate({ scrollTop: $target.offset().top - 72 }, 320);
  });

  $('[href^="mailto:"]').on('click', function () {
    $(this).attr('data-clicked-at', new Date().toISOString());
  });
});
