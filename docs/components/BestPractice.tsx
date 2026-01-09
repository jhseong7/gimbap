import React from 'react';
import styles from './BestPractice.module.css';

interface BestPracticeProps {
  doList?: string[];
  dontList?: string[];
  children?: React.ReactNode;
}

export function BestPractice({ doList, dontList, children }: BestPracticeProps) {
  return (
    <div className={styles.container}>
      {children && <div className={styles.intro}>{children}</div>}

      <div className={styles.grid}>
        {doList && doList.length > 0 && (
          <div className={styles.column}>
            <div className={`${styles.header} ${styles.do}`}>
              <span className={styles.icon}>✅</span>
              <strong>Do</strong>
            </div>
            <ul className={styles.list}>
              {doList.map((item, idx) => (
                <li key={idx}>{item}</li>
              ))}
            </ul>
          </div>
        )}

        {dontList && dontList.length > 0 && (
          <div className={styles.column}>
            <div className={`${styles.header} ${styles.dont}`}>
              <span className={styles.icon}>❌</span>
              <strong>Don't</strong>
            </div>
            <ul className={styles.list}>
              {dontList.map((item, idx) => (
                <li key={idx}>{item}</li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}

export default BestPractice;
