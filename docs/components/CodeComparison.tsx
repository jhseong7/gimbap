import React from 'react';
import styles from './CodeComparison.module.css';

interface CodeComparisonProps {
  before: string;
  after: string;
  beforeLabel?: string;
  afterLabel?: string;
  language?: string;
}

export function CodeComparison({
  before,
  after,
  beforeLabel = 'Before',
  afterLabel = 'After',
  language = 'go'
}: CodeComparisonProps) {
  return (
    <div className={styles.comparison}>
      <div className={styles.side}>
        <div className={styles.label}>
          <span className={styles.icon}>❌</span>
          {beforeLabel}
        </div>
        <pre className={styles.codeBlock}>
          <code className={`language-${language}`}>{before}</code>
        </pre>
      </div>
      <div className={styles.arrow}>→</div>
      <div className={styles.side}>
        <div className={styles.label}>
          <span className={styles.icon}>✅</span>
          {afterLabel}
        </div>
        <pre className={styles.codeBlock}>
          <code className={`language-${language}`}>{after}</code>
        </pre>
      </div>
    </div>
  );
}

export default CodeComparison;
