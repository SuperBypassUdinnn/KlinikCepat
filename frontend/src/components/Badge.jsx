import './Badge.css';

/**
 * Badge — label status triage atau antrean.
 *
 * @param {'merah'|'kuning'|'hijau'|'menunggu'|'selesai'|'dilewati'} status
 * @param {'sm'|'lg'} size
 */
export default function Badge({ status, size = 'sm', children, className = '' }) {
  const statusLower = status?.toLowerCase() || '';

  const label = children || status;

  const classes = [
    'badge',
    `badge-${statusLower}`,
    size === 'lg' && 'badge-lg',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <span className={classes}>
      <span className="badge-dot" />
      {label}
    </span>
  );
}
