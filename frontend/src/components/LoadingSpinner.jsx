import './LoadingSpinner.css';

/**
 * LoadingSpinner — indikator loading.
 *
 * @param {'sm'|'md'|'lg'} size
 * @param {boolean} fullPage — tampilkan overlay full-page
 * @param {string} text — teks di bawah spinner
 */
export default function LoadingSpinner({ size = 'md', fullPage = false, text }) {
  const spinnerClass = [
    'spinner-ring',
    size === 'sm' && 'spinner-sm',
    size === 'lg' && 'spinner-lg',
  ]
    .filter(Boolean)
    .join(' ');

  const wrapperClass = [
    'spinner-overlay',
    fullPage && 'spinner-fullpage',
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div className={wrapperClass}>
      <div className={spinnerClass} />
      {text && <p className="spinner-text">{text}</p>}
    </div>
  );
}
