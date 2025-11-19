import React from 'react';
import styles from './Steps.module.css';

interface StepsProps {
  children: React.ReactNode;
}

interface StepProps {
  title: string;
  children: React.ReactNode;
}

export function Steps({ children }: StepsProps) {
  return <div className={styles.steps}>{children}</div>;
}

export function Step({ title, children }: StepProps) {
  return (
    <div className={styles.step}>
      <div className={styles.stepTitle}>
        <div className={styles.stepNumber}></div>
        <h3>{title}</h3>
      </div>
      <div className={styles.stepContent}>{children}</div>
    </div>
  );
}

export default Steps;
