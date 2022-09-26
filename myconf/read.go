package myconf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"time"
)

func WatchFile(fileName string, paths []string, conf interface{}, delayTime time.Duration) (*Watch, error) {
	var v = viper.New()
	v.SetConfigName(fileName)
	v.SetConfigType("yaml")
	for _, path := range paths {
		v.AddConfigPath(path)
	}

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = v.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	w := NewWatch(delayTime)

	v.OnConfigChange(func(in fsnotify.Event) {
		err := v.Unmarshal(conf)
		if err != nil {
			log.Printf("[error] unable to unmarshal: %v\n", err)
		} else {
			log.Println("watch file reload: ", in.String())
			w.Send()
		}
	})
	v.WatchConfig()
	return w, nil
}

func ReadFile(fileName string, paths []string, conf interface{}) error {
	var v = viper.New()
	v.SetConfigName(fileName)
	v.SetConfigType("yaml")
	for _, path := range paths {
		v.AddConfigPath(path)
	}
	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	v.Unmarshal(conf)

	w := new(Watch)

	v.OnConfigChange(func(e fsnotify.Event) {
		err := v.Unmarshal(conf)
		log.Printf("[error] unable to unmarshal: %v\n", err)
		w.Send()
	})

	viper.WatchConfig()

	return nil
}

func ReadRemote(provider string, endpoint string, path string, conf interface{}) error {
	var v = viper.New()

	v.AddRemoteProvider(provider, endpoint, path)
	v.SetConfigType("yaml")
	err := v.ReadRemoteConfig()
	if err != nil {
		return err
	}
	err = v.Unmarshal(conf)
	if err != nil {
		log.Printf("[error] unable to unmarshal: %v\n", err)
		return err
	}

	go func() {
		for {
			time.Sleep(time.Second * 5)
			err = v.WatchRemoteConfig()

			if err != nil {
				log.Printf("[error] unable to read remote config: %v\n", err)
				continue
			}

			v.Unmarshal(conf)
		}
	}()

	return nil
}

func WatchReadRemote(provider string, endpoint string, path string, conf interface{}) (*Watch, error) {
	var v = viper.New()

	err := v.AddRemoteProvider(provider, endpoint, path)
	if err != nil {
		return nil, err
	}

	v.SetConfigType("yaml")
	err = v.ReadRemoteConfig()
	if err != nil {
		log.Printf("[error] unable to read remote config: %v\n", err)
		return nil, err
	}
	err = v.Unmarshal(conf)
	if err != nil {
		log.Printf("[error] unable to unmarshal: %v\n", err)
		return nil, err
	}

	w := new(Watch)

	go func() {
		for {
			time.Sleep(time.Second * 5)

			err := v.WatchRemoteConfig()

			if err != nil {
				log.Printf("[error] unable to read remote config: %v\n", err)
				continue
			}

			err = v.Unmarshal(conf)
			if err != nil {
				log.Printf("[error] unable to unmarshal: %v\n", err)
				continue
			}

			w.Send()
		}
	}()

	return w, nil
}
