import React from 'react';
import styles from './ApiCard.module.css';

interface ApiCardProps {
  name: string;
  signature: string;
  description: string;
  parameters?: Array<{
    name: string;
    type: string;
    description: string;
    optional?: boolean;
  }>;
  returns?: {
    type: string;
    description: string;
  };
  example?: string;
}

export function ApiCard({ name, signature, description, parameters, returns, example }: ApiCardProps) {
  return (
    <div className={styles.card}>
      <div className={styles.header}>
        <h3 className={styles.name}>{name}</h3>
      </div>

      <div className={styles.signature}>
        <code>{signature}</code>
      </div>

      <div className={styles.description}>
        {description}
      </div>

      {parameters && parameters.length > 0 && (
        <div className={styles.section}>
          <h4>Parameters</h4>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Name</th>
                <th>Type</th>
                <th>Description</th>
              </tr>
            </thead>
            <tbody>
              {parameters.map((param, idx) => (
                <tr key={idx}>
                  <td>
                    <code>{param.name}</code>
                    {param.optional && <span className={styles.optional}> (optional)</span>}
                  </td>
                  <td><code>{param.type}</code></td>
                  <td>{param.description}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {returns && (
        <div className={styles.section}>
          <h4>Returns</h4>
          <p><code>{returns.type}</code> - {returns.description}</p>
        </div>
      )}

      {example && (
        <div className={styles.section}>
          <h4>Example</h4>
          <pre className={styles.example}>
            <code>{example}</code>
          </pre>
        </div>
      )}
    </div>
  );
}

export default ApiCard;
