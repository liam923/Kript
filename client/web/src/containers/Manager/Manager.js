import React, { Component } from 'react';

import Datum from '../../components/Datum/Datum';
import NavigationBar from '../../components/common/NavigationBar/NavigationBar';

import styles from './Manager.module.css';

class Manager extends Component {
    state = {
        info: {
            passwords: [
                { url: 'www.amazon.com', username: 'abc@gmail.com', password: 'password123' },
                { url: 'www.fb.com', username: 'abc@gmail.com', password: 'password123' },
                { url: 'www.gmail.com', username: 'abc@gmail.com', password: 'jaimeisfoxy' },
            ],
            creditcards: [
                { number: '4111 1111 1111 1111', name: 'Andy B. Calvin', expMonth: '02', expYear: '21', cvv: '123', desc: 'Visa' },
                { number: '4111 1111 1111 4444', name: 'Andy B. Calvin', expMonth: '04', expYear: '25', cvv: '123', desc: 'Mastercard' },
            ],
            codes: [
                { code: '1234', type: 'Bank Account', desc: 'Bank Account' }
            ],
            notes: [],
        },
        filter: ''
    }

    render() {
        let items = Object.keys(this.state.info)
            .map(pKey => {
                return this.state.info[pKey].map((el, i) => {
                    return <Datum key={pKey + i} type={pKey} datum={el} />
                });
            }).reduce((arr, el) => {
                return arr.concat(el);
            }, []);

        console.log(items);

        return (
            <div style={{display: 'flex'}}>
                <div style={{width: '200px'}}>
                    <div id='spacer' style={{height: '118px'}}></div>
                    <NavigationBar />
                </div>
                <div style={{marginLeft: '32px', flexGrow: '1'}}>
                    <h1>Kript</h1>
                    <h2>My Data</h2>
                    <div className={styles.Manager}>
                        {items}
                    </div>
                </div>

            </div>
        );
    }
}

export default Manager;