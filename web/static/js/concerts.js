// Просмотр фото выступлений на весь экран.
// Без JS обложка просто открывает первое фото в полном размере.
(() => {
  const dialog = document.getElementById('lightbox');
  if (!dialog || typeof dialog.showModal !== 'function') return;

  const img = dialog.querySelector('.lightbox__img');
  const caption = dialog.querySelector('.lightbox__caption');
  const prevBtn = dialog.querySelector('.lightbox__prev');
  const nextBtn = dialog.querySelector('.lightbox__next');

  let photos = [];   // ссылки на фото текущей галереи
  let index = 0;
  let title = '';

  const show = (i) => {
    index = (i + photos.length) % photos.length;
    img.src = photos[index].href;
    img.alt = `${title}, фото ${index + 1} из ${photos.length}`;
    caption.textContent = `${title}: ${index + 1} из ${photos.length}`;

    // Подгружаем соседние фото заранее
    [index + 1, index - 1].forEach((j) => {
      new Image().src = photos[(j + photos.length) % photos.length].href;
    });

    prevBtn.hidden = nextBtn.hidden = photos.length < 2;
  };

  // Клик по обложке: собираем все ссылки этой галереи (в том числе скрытые)
  document.addEventListener('click', (e) => {
    const link = e.target.closest('a[data-gallery]');
    if (!link) return;
    e.preventDefault();

    photos = [...document.querySelectorAll(`a[data-gallery="${CSS.escape(link.dataset.gallery)}"]`)];
    title = link.closest('[data-title]')?.dataset.title ?? '';
    show(photos.indexOf(link));
    dialog.showModal();
  });

  prevBtn.addEventListener('click', () => show(index - 1));
  nextBtn.addEventListener('click', () => show(index + 1));
  dialog.querySelector('.lightbox__close').addEventListener('click', () => dialog.close());

  // Клик по тёмному фону закрывает просмотр
  dialog.addEventListener('click', (e) => {
    if (!e.target.closest('.lightbox__btn, .lightbox__img, .lightbox__caption')) dialog.close();
  });

  // Стрелки на клавиатуре (Escape закрывает диалог сам)
  dialog.addEventListener('keydown', (e) => {
    if (e.key === 'ArrowLeft') show(index - 1);
    if (e.key === 'ArrowRight') show(index + 1);
  });

  // Свайп на телефоне
  let startX = null;
  dialog.addEventListener('touchstart', (e) => {
    startX = e.touches[0].clientX;
  }, { passive: true });
  dialog.addEventListener('touchend', (e) => {
    if (startX === null) return;
    const dx = e.changedTouches[0].clientX - startX;
    startX = null;
    if (Math.abs(dx) > 50 && photos.length > 1) show(index + (dx < 0 ? 1 : -1));
  });

  dialog.addEventListener('close', () => img.removeAttribute('src'));
})();
