import React from 'react';

import AllInclusiveIcon from '@material-ui/icons/AllInclusive';
import VpnKeyIcon from '@material-ui/icons/VpnKey';
import CreditCardIcon from '@material-ui/icons/CreditCard';
import NoteIcon from '@material-ui/icons/Note';
import FiberPinIcon from '@material-ui/icons/FiberPin';

import styles from './NavigationItem.module.css';

const navigationItem = (props) => {
    let output = null;

    switch (props.navType) {
        case ('all'):
            output = <span onClick={props.clicked} className={props.active ? styles.active : null}><AllInclusiveIcon className={styles.icon}/><span className={styles.NavText}>All</span></span>
            break;
        case ('passwords'):
            output = <span onClick={props.clicked} className={props.active ? styles.active : null}><VpnKeyIcon className={styles.icon}/><span>Passwords</span></span>
            break;
        case ('creditcards'):
            output = <span onClick={props.clicked} className={props.active ? styles.active : null}><CreditCardIcon className={styles.icon}/><span>Credit Cards</span></span>
            break;
        case ('notes'):
            output = <span onClick={props.clicked} className={props.active ? styles.active : null}><NoteIcon className={styles.icon}/><span>Notes</span></span>
            break;
        case ('codes'):
            output = <span onClick={props.clicked} className={props.active ? styles.active : null}><FiberPinIcon className={styles.icon}/><span>Codes</span></span>
            break;
        default:
            output = null;
            break;
    }

    return (
        output
    )
};

export default navigationItem;