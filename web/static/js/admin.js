// Админские страницы: список песен, превью обложки, загрузка фото, подтверждение удаления.
// Каждый блок работает, только если нужные элементы есть на странице.
(() => {
  // ---------- Подтверждение удаления ----------
  document.querySelectorAll('form[data-confirm]').forEach((form) => {
    form.addEventListener('submit', (e) => {
      if (!confirm(form.dataset.confirm)) e.preventDefault();
    });
  });

  // ---------- Обложка альбома: превью выбранного файла ----------
  const coverInput = document.getElementById('cover');
  const coverPreview = document.getElementById('cover-preview');
  const coverEmpty = document.getElementById('cover-empty');
  if (coverInput && coverPreview) {
    coverInput.addEventListener('change', () => {
      const file = coverInput.files[0];
      if (!file) return;
      coverPreview.src = URL.createObjectURL(file);
      coverPreview.hidden = false;
      if (coverEmpty) coverEmpty.hidden = true;
    });
  }

  // ---------- Список песен ----------
  // Поля песен называются songs[0].title, songs[0].duration, songs[1].title, ...
  // Индексы должны идти подряд в порядке строк, поэтому после любого изменения списка их пересчитываем.
  const songs = document.getElementById('songs');
  if (songs) {
    const template = document.getElementById('song-row-template');
    const total = document.getElementById('songs-total');

    const toSeconds = (value) => {
      const m = /^(\d{1,3}):([0-5]\d)$/.exec(value.trim());
      return m ? Number(m[1]) * 60 + Number(m[2]) : 0;
    };
    const pad = (n) => String(n).padStart(2, '0');
    const format = (sec) => `${pad(Math.floor(sec / 60))}:${pad(sec % 60)}`;

    const reindex = () => {
      [...songs.children].forEach((row, i) => {
        row.querySelectorAll('input[name^="songs["]').forEach((input) => {
          input.name = input.name.replace(/^songs\[\d*\]/, `songs[${i}]`);
        });
      });
    };

    const refresh = () => {
      reindex();
      const rows = [...songs.children];
      const sec = rows.reduce(
        (sum, row) => sum + toSeconds(row.querySelector('input[name$=".duration"]').value), 0);
      total.textContent = `Песен: ${rows.length}, общая длительность: ${format(sec)}`;
    };

    const addRow = (focus) => {
      songs.append(template.content.cloneNode(true));
      refresh();
      if (focus) songs.lastElementChild.querySelector('input').focus();
    };

    document.getElementById('add-song').addEventListener('click', () => addRow(true));

    songs.addEventListener('click', (e) => {
      const row = e.target.closest('.songs__row');
      if (!row) return;
      if (e.target.closest('[data-song-up]') && row.previousElementSibling) {
        row.previousElementSibling.before(row);
        refresh();
      } else if (e.target.closest('[data-song-down]') && row.nextElementSibling) {
        row.nextElementSibling.after(row);
        refresh();
      } else if (e.target.closest('[data-song-remove]')) {
        row.remove();
        refresh();
      }
    });

    songs.addEventListener('input', refresh);

    // На всякий случай пересчитываем индексы прямо перед отправкой формы
    songs.closest('form').addEventListener('submit', reindex);

    refresh();
  }

  // ---------- Фото выступления: накопление файлов и предпросмотр ----------
  const photosInput = document.getElementById('photos');
  if (photosInput) {
    const list = document.getElementById('new-photos');
    const drop = photosInput.closest('.uploader').querySelector('.uploader__drop');
    let files = [];

    const tile = (file, index) => {
      const li = document.createElement('li');
      li.className = 'photo';

      const img = new Image();
      img.className = 'photo__img';
      img.alt = '';
      img.src = URL.createObjectURL(file);
      img.addEventListener('load', () => URL.revokeObjectURL(img.src));

      const name = document.createElement('span');
      name.className = 'photo__name';
      name.textContent = file.name;

      const remove = document.createElement('button');
      remove.type = 'button';
      remove.className = 'btn btn--edit';
      remove.textContent = 'Убрать';
      remove.setAttribute('aria-label', `Убрать файл ${file.name}`);
      remove.addEventListener('click', () => {
        files.splice(index, 1);
        sync();
      });

      li.append(img, name, remove);
      return li;
    };

    // Поле выбора файлов всегда хранит все накопленные файлы, чтобы они ушли на сервер вместе с формой
    const sync = () => {
      const dt = new DataTransfer();
      files.forEach((f) => dt.items.add(f));
      photosInput.files = dt.files;
      list.replaceChildren(...files.map(tile));
    };

    const add = (picked) => {
      const fresh = [...picked].filter((f) => f.type.startsWith('image/')
        && !files.some((g) => g.name === f.name && g.size === f.size && g.lastModified === f.lastModified));
      files = files.concat(fresh);
      sync();
    };

    photosInput.addEventListener('change', () => add(photosInput.files));

    ['dragenter', 'dragover'].forEach((type) => drop.addEventListener(type, (e) => {
      e.preventDefault();
      drop.classList.add('is-over');
    }));
    ['dragleave', 'drop'].forEach((type) => drop.addEventListener(type, (e) => {
      e.preventDefault();
      drop.classList.remove('is-over');
    }));
    drop.addEventListener('drop', (e) => add(e.dataTransfer.files));
  }
})();
