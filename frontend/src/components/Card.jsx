import './Card.css';

/**
 * Card — container card reusable.
 *
 * @param {'default'|'glass'|'flat'} variant
 * @param {'merah'|'kuning'|'hijau'|'primary'} accent — warna border atas
 */
export default function Card({
  children,
  variant = 'default',
  accent,
  className = '',
  ...props
}) {
  const classes = [
    'card',
    variant === 'glass' && 'card-glass',
    variant === 'flat' && 'card-flat',
    accent && `card-accent-${accent}`,
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div className={classes} {...props}>
      {children}
    </div>
  );
}
