import React from 'react';

import AllInclusiveIcon from '@material-ui/icons/AllInclusive';
import VpnKeyIcon from '@material-ui/icons/VpnKey';
import CreditCardIcon from '@material-ui/icons/CreditCard';
import NoteIcon from '@material-ui/icons/Note';
import FiberPinIcon from '@material-ui/icons/FiberPin';
import Button from '../Button/Button';

import styles from './NavigationBar.module.css';

const navigationBar = (props) => (
    <div className={styles.NavigationBar}>
        <div>
            <a href='/' className={styles.active}><AllInclusiveIcon className={styles.icon}/><span className={styles.NavText}>All</span></a>
            <a href='/'><VpnKeyIcon className={styles.icon}/><span>Passwords</span></a>
            <a href='/'><CreditCardIcon className={styles.icon}/><span>Credit Cards</span></a>
            <a href='/'><NoteIcon className={styles.icon}/><span>Notes</span></a>
            <a href='/'><FiberPinIcon className={styles.icon}/><span>Codes</span></a>
        </div>
        <Button>Add Data</Button>
        <Button>Log Out</Button>
    </div>
);

export default navigationBar;