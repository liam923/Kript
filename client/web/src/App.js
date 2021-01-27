import React from 'react';

import styles from './App.module.css';

import Layout from './hoc/Layout/Layout';
import Manager from './containers/Manager/Manager';

function App() {
  return (
    <div className={styles.App}>
        <Layout>
            <Manager />
        </Layout>
    </div>
  );
}

export default App;
