import React from 'react';
import styles from './FeatureCard.module.css';

interface FeatureCardProps {
  icon?: string;
  title: string;
  description: string;
  link?: string;
}

export function FeatureCard({ icon, title, description, link }: FeatureCardProps) {
  const content = (
    <>
      {icon && <div className={styles.icon}>{icon}</div>}
      <h3 className={styles.title}>{title}</h3>
      <p className={styles.description}>{description}</p>
    </>
  );

  if (link) {
    return (
      <a href={link} className={styles.card}>
        {content}
        <div className={styles.arrow}>→</div>
      </a>
    );
  }

  return <div className={styles.card}>{content}</div>;
}

interface FeatureGridProps {
  children: React.ReactNode;
}

export function FeatureGrid({ children }: FeatureGridProps) {
  return <div className={styles.grid}>{children}</div>;
}

export default FeatureCard;
