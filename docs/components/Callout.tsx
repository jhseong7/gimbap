import React from 'react';
import styles from './Callout.module.css';

export type CalloutType = 'info' | 'warning' | 'error' | 'success' | 'tip';

interface CalloutProps {
  type?: CalloutType;
  title?: string;
  children: React.ReactNode;
}

const icons = {
  info: '💡',
  warning: '⚠️',
  error: '❌',
  success: '✅',
  tip: '💡',
};

const titles = {
  info: 'Info',
  warning: 'Warning',
  error: 'Error',
  success: 'Success',
  tip: 'Tip',
};

export function Callout({ type = 'info', title, children }: CalloutProps) {
  const displayTitle = title || titles[type];

  return (
    <div className={`${styles.callout} ${styles[type]}`}>
      <div className={styles.title}>
        <span className={styles.icon}>{icons[type]}</span>
        <strong>{displayTitle}</strong>
      </div>
      <div className={styles.content}>{children}</div>
    </div>
  );
}

export default Callout;
