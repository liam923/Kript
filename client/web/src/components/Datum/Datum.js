import React from 'react';

import Aux from '../../hoc/Aux/Aux';

import styles from './Datum.module.css';

const datum = (props) => {
    let datum = null;

    switch (props.type) {
        case ('passwords'):
            datum = (
                <Aux>
                    <h4>{props.datum.url}</h4>
                    <p>{props.datum.username}</p>
                </Aux>
            );
            break;
        case ('creditcards'):
            datum = (
                <Aux>
                    <h4>{props.datum.desc}</h4>
                    <p>Ending with {props.datum.number.slice(props.datum.number.length - 4)}</p>
                </Aux>
            );
            break;
        case ('notes'):
            datum = (
                <Aux>
                    <h4>{props.datum.title}</h4>
                </Aux>
            );
            break;
        case ('codes'):
            datum = (
                <Aux>
                    <h4>{props.datum.desc}</h4>
                    <p>{props.datum.type}</p>
                </Aux>
            );
            break;
        default:
            datum = null;
            break;
    }

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