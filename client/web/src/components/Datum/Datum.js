import React from 'react';

import Aux from '../../hoc/Aux/Aux';

import styles from './Datum.module.css';
import {DatumType} from './DatumType/DatumType';
import {Password} from './DatumType/Password';

const datum = (props) => {
    let datum = DatumType.getType(props.type).card(props);

    return (
        // revise onClick method; should it be ID based or not?
        <div className={styles.Datum} onClick={() => props.clicked(props.datum)}>
            <div style={{ backgroundColor: 'beige', height: '100px', borderRadius: '15px 15px 0 0' }}></div>
            <div style={{ padding: '20px' }}>
                {datum}
            </div>
        </div>
    );
}

export default datum;