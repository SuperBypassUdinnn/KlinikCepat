import './Button.css';

/**
 * Button — komponen tombol reusable dengan variasi gaya dan ukuran.
 *
 * @param {'primary'|'secondary'|'danger'|'success'|'warning'|'ghost'|'outline'} variant
 * @param {'sm'|'md'|'lg'} size
 * @param {boolean} loading — tampilkan spinner & disable
 * @param {boolean} block — full width
 */
export default function Button({
  children,
  variant = 'primary',
  size = 'md',
  loading = false,
  block = false,
  className = '',
  disabled,
  ...props
}) {
  const classes = [
    'btn',
    `btn-${variant}`,
    `btn-${size}`,
    block && 'btn-block',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <button className={classes} disabled={disabled || loading} {...props}>
      {loading && <span className="btn-spinner" />}
      {children}
    </button>
  );
}
