/* app.js — handles search form and rendering the profile API response */

const form    = document.getElementById('search-form');
const input   = document.getElementById('username-input');
const btn     = document.getElementById('search-btn');

const errBanner  = document.getElementById('error-banner');
const errMsg     = document.getElementById('error-message');
const loading    = document.getElementById('loading-state');
const result     = document.getElementById('profile-result');

// --- helpers ----------------------------------------------------------------

function show(el) { el.classList.remove('hidden'); }
function hide(el) { el.classList.add('hidden'); }

function showError(msg) {
  errMsg.textContent = msg;
  show(errBanner);
  hide(loading);
  hide(result);
}

function scoreClass(n) {
  if (n >= 70) return 'good';
  if (n >= 40) return 'ok';
  return 'low';
}

function setBar(barId, numId, value) {
  const bar = document.getElementById(barId);
  const num = document.getElementById(numId);
  bar.style.width = value + '%';
  bar.className   = 'score-bar ' + scoreClass(value);
  num.textContent = value;
}

function avatarUrl(login) {
  return 'https://github.com/' + encodeURIComponent(login) + '.png?size=128';
}

// --- render -----------------------------------------------------------------

function renderProfile(data) {
  const { user, repositories, analysis } = data;

  // Profile header
  const avatar = document.getElementById('profile-avatar');
  avatar.src = avatarUrl(user.login);
  avatar.alt = user.login + ' avatar';

  document.getElementById('profile-name').textContent    = user.name || user.login;
  document.getElementById('profile-login').textContent   = '@' + user.login;
  document.getElementById('profile-bio').textContent     = user.bio || '';
  document.getElementById('stat-repos').textContent      = user.public_repos;
  document.getElementById('stat-followers').textContent  = user.followers;
  document.getElementById('stat-following').textContent  = user.following;

  // Scores
  const overall = analysis.overall_score;
  document.getElementById('overall-badge').textContent = 'Overall ' + overall + '/100';

  setBar('bar-doc',      'num-doc',      analysis.documentation_score);
  setBar('bar-repo',     'num-repo',     analysis.repository_score);
  setBar('bar-activity', 'num-activity', analysis.activity_score);

  // Recommendations
  const recList = document.getElementById('recommendation-list');
  recList.innerHTML = '';
  (analysis.recommendations || []).forEach(function(rec) {
    const li = document.createElement('li');
    li.textContent = rec;
    recList.appendChild(li);
  });

  // Repositories
  const repoList = document.getElementById('repo-list');
  repoList.innerHTML = '';

  document.getElementById('repo-count').textContent =
    repositories ? repositories.length + ' repos' : '';

  (repositories || []).forEach(function(repo) {
    const li = document.createElement('li');
    li.className = 'repo-item';

    const badges = [
      repo.is_fork     ? '<span class="badge badge-fork">Fork</span>'         : '',
      repo.archived    ? '<span class="badge badge-archived">Archived</span>' : '',
    ].join('');

    const stars = repo.stars
      ? '<span class="repo-meta-item">' +
          starIcon() + ' ' + repo.stars +
        '</span>'
      : '';

    const lang = repo.language
      ? '<span class="repo-meta-item">' + esc(repo.language) + '</span>'
      : '';

    li.innerHTML =
      '<div class="repo-item-header">' +
        '<a class="repo-name" href="' + esc(repo.html_url) + '" target="_blank" rel="noopener">' +
          esc(repo.name) +
        '</a>' +
        '<div class="repo-badges">' + badges + '</div>' +
      '</div>' +
      (repo.description
        ? '<p class="repo-description">' + esc(repo.description) + '</p>'
        : '') +
      '<div class="repo-meta">' + stars + lang + '</div>';

    repoList.appendChild(li);
  });

  show(result);
}

function starIcon() {
  return '<svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg">' +
    '<path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>' +
    '</svg>';
}

// simple HTML escape to avoid XSS from API data
function esc(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

// --- fetch ------------------------------------------------------------------

async function fetchProfile(username) {
  hide(errBanner);
  hide(result);
  show(loading);
  btn.disabled = true;

  try {
    const res = await fetch('/profile/' + encodeURIComponent(username));
    const body = await res.json();

    if (!res.ok) {
      const msg = body.error || 'Request failed (' + res.status + ')';
      showError(msg);
      return;
    }

    hide(loading);
    renderProfile(body);
  } catch (err) {
    showError('Could not reach the server. Is it running?');
  } finally {
    btn.disabled = false;
  }
}

// --- events -----------------------------------------------------------------

form.addEventListener('submit', function(e) {
  e.preventDefault();
  const username = input.value.trim();
  if (!username) return;
  fetchProfile(username);
});

// Support ?u=username in the URL for shareable links
(function() {
  const params = new URLSearchParams(window.location.search);
  const u = params.get('u');
  if (u) {
    input.value = u;
    fetchProfile(u);
  }
})();
